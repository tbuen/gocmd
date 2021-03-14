// Package config handles the persistent configuration data.
package config

var (
	bookmarks *Bookmarks
)

// Load reads the configuration data from disk.
func Load() {
	readApps(filenameApps)

	/*buf, err = load(filenameBookmarks)
	if err == nil {
		convertBookmarks(read, buf, &bookmarks)
	}*/

	readTabs(filenameTabs)
}

// Save writes the configuration data to disk.
func Save() {
	writeTabs(filenameTabs)

	//cfg.bookmarks.Save(filenameBookmarks)
}

func Bookmarksi() *Bookmarks {
	return bookmarks
}
