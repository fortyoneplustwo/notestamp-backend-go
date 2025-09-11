package cleanup

type MetadataRemover interface {
	MetadataRemove(uid string, id string) error
}

type MediaRemover interface {
	MediaRemove(uid string, oid string) error
}

type NotesRemover interface {
	NotesRemove(uid string, oid string) error
}
