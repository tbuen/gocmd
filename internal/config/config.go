// Package config handles the persistent configuration data.
package config

// Load reads the configuration data from disk.
func Load() {
	readApps(filenameApps)
	readBookmarks(filenameBookmarks)
	readTabs(filenameTabs)
}

// Save writes the configuration data to disk.
func Save() {
	writeBookmarks(filenameBookmarks)
	writeTabs(filenameTabs)
}
