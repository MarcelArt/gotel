package web

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v3"
)

// LicenseFile represents a single license or notice file.
type LicenseFile struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// LicenseInfo represents the license(s) associated with a third-party package.
type LicenseInfo struct {
	Package string        `json:"package"`
	Files   []LicenseFile `json:"files"`
}

// LicensesViewModel represents the view parameters for the licenses page.
type LicensesViewModel struct {
	BaseViewModel
	Licenses []LicenseInfo
	Error    string
}

var (
	licensesCache  []LicenseInfo
	licensesLoaded bool
	licensesMu     sync.RWMutex
)

// findLicensesDir searches for the THIRD_PARTY_LICENSES directory in common paths.
func findLicensesDir() (string, bool) {
	if dir := os.Getenv("LICENSES_DIR"); dir != "" {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return dir, true
		}
	}
	paths := []string{
		"THIRD_PARTY_LICENSES",
		"../THIRD_PARTY_LICENSES",
		"../../THIRD_PARTY_LICENSES",
		"/app/THIRD_PARTY_LICENSES",
	}
	for _, p := range paths {
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return p, true
		}
	}
	return "", false
}

// isLicenseFile returns true if the filename matches a license or notice pattern.
func isLicenseFile(name string) bool {
	name = strings.ToLower(name)
	// Check exact matches
	if name == "license" || name == "licence" || name == "notice" || name == "copying" {
		return true
	}
	// Check prefix with extensions
	if strings.HasPrefix(name, "license.") || strings.HasPrefix(name, "licence.") || strings.HasPrefix(name, "notice.") || strings.HasPrefix(name, "copying.") {
		ext := filepath.Ext(name)
		return ext == ".txt" || ext == ".md"
	}
	return false
}

// LoadLicenses walks the given directory and aggregates license files by package.
func LoadLicenses(dir string) ([]LicenseInfo, error) {
	licenseMap := make(map[string][]LicenseFile)

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		filename := d.Name()
		if isLicenseFile(filename) {
			rel, err := filepath.Rel(dir, path)
			if err != nil {
				return nil
			}

			pkg := filepath.Dir(rel)
			if pkg == "." {
				pkg = "gotel"
			}

			contentBytes, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			content := string(contentBytes)
			content = strings.ReplaceAll(content, "\r\n", "\n")

			licenseMap[pkg] = append(licenseMap[pkg], LicenseFile{
				Name:    filename,
				Content: content,
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	var list []LicenseInfo
	for pkg, files := range licenseMap {
		list = append(list, LicenseInfo{
			Package: pkg,
			Files:   files,
		})
	}

	// Sort package names alphabetically
	sort.Slice(list, func(i, j int) bool {
		return list[i].Package < list[j].Package
	})

	return list, nil
}

// getLicenses returns the cached license list or parses them if not already cached.
func getLicenses() ([]LicenseInfo, error) {
	licensesMu.RLock()
	if licensesLoaded {
		res := licensesCache
		licensesMu.RUnlock()
		return res, nil
	}
	licensesMu.RUnlock()

	licensesMu.Lock()
	defer licensesMu.Unlock()

	if licensesLoaded {
		return licensesCache, nil
	}

	dir, found := findLicensesDir()
	if !found {
		return nil, os.ErrNotExist
	}

	list, err := LoadLicenses(dir)
	if err != nil {
		return nil, err
	}

	licensesCache = list
	licensesLoaded = true
	return licensesCache, nil
}

// LicensesGet handles GET /licenses requests and displays the attributions tab.
func (h *WebHandler) LicensesGet(c fiber.Ctx) error {
	userID := c.Locals("userId")
	if userID == nil {
		return c.Redirect().To("/login")
	}

	user, err := h.userService.GetByID(c, userID)
	if err != nil {
		return h.LogoutPost(c)
	}

	licenses, err := getLicenses()
	var errMsg string
	if err != nil {
		errMsg = "Third-party licenses directory ('THIRD_PARTY_LICENSES') was not found or could not be walked. Please ensure you have run 'go-licenses save THIRD_PARTY_LICENSES' to generate the attributions directory."
	}

	vm := LicensesViewModel{
		BaseViewModel: BaseViewModel{
			Title:       "Attribution Licenses - Gotel",
			ActiveTab:   "licenses",
			User:        user,
			Permissions: getPermissions(c),
		},
		Licenses: licenses,
		Error:    errMsg,
	}

	return h.renderTab(c, "licenses", vm)
}
