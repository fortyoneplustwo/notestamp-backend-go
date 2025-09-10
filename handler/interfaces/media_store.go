package handler

type MediaRemover interface {
	MediaRemove(uid string, oid string) error
	MediaRemoveAll(uid string) error
}

type MediaUploadPathGetter interface {
	MediaGetUploadPath(uid string, oid string) string
}

type MediaChecker interface {
	MediaCheck(uid string, oid string) (bool, error)
}

type MediaCheckRemover interface {
	MediaCheck(uid string, oid string) (bool, error)
	MediaRemove(uid string, oid string) error
	MediaRemoveAll(uid string) error
}
