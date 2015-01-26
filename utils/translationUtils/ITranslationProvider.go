package translationUtils

type ITranslationProvider interface {
	GetTranslation(translationKey string, params map[string]string) string
}
