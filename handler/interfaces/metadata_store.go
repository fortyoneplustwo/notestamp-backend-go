package handler

import "notestamp/metadata"

type MetadataGetter interface {
	MetadataGet(uid string, oid string) (metadata.Metadata, error)
}

type MetadataLister interface {
	MetadataList(uid string) ([]metadata.Metadata, error)
}

type MetadataAdder interface {
	MetadataAdd(uid string, m metadata.Metadata) error
}

type MetadataRemover interface {
	MetadataRemove(uid string, id string) error
	MetadataRemoveAll(uid string) error
}

type MetadataUpdater interface {
	MetadataUpdate(uid string, m metadata.Metadata) error
}

type MetadataGetRemover interface {
	MetadataGet(uid string, id string) (metadata.Metadata, error)
	MetadataRemove(uid string, id string) error
	MetadataRemoveAll(uid string) error
}

