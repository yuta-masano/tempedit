// Package tempedit provides a temporary file base user input interface using an external editor.
//
// Here is a simple example.
//
//	tmpFile, err := tempedit.NewTempFile("", ".tmp.")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer tempedit.Clean(tmpFile) // It is the caller's responsibility to remove the file when no longer needed.
//
//	// You can write some contents into the temporary file.
//	if err := tmpFile.Write("hello 世界\n"); err != nil {
//		log.Fatal(err)
//	}
//	// You can let user write contents into the temporary file using an external editor.
//	vi := tempedit.NewEditor("vi")
//	if err := tmpFile.OpenWith(vi); err != nil {
//		log.Fatal(err)
//	}
//	if changed, err := tmpFile.IsChanged(); !changed {
//		log.Fatal(err)
//	}
//
//	// You can get the contents of the tempfile.
//	userInputed := tmpFile.String()
package tempedit
