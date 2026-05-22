package name

// Formatter интерфейс для форматирования имени. В домене логика форматирования отсутствует.
// Все реализации данного интерфейса находятся в слое презентации.
type Formatter interface {
	Format(name Name) string
}
