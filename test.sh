cp -pr ./sample/* ./oldfiles/
go build && ./cleaner ./oldfiles/ '.*' 5

