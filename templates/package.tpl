<!DOCTYPE html>
<html>
    <head>
        <meta name="go-import" content="{{ .CanonicalURL }} git https://{{ .Repo }}">
        <meta name="go-source" content="{{ .CanonicalURL }} https://{{ .Repo }} https://{{ .Repo }}/tree/master{/dir} https://{{ .Repo }}/tree/master{/dir}/{file}#L{line}">
        <meta http-equiv="refresh" content="0; url=https://godoc.org/{{ .CanonicalURL }}">
    </head>
    <body>
        Nothing to see here. Please <a href="https://godoc.org/{{ .CanonicalURL }}">move along</a>.
    </body>
</html>
