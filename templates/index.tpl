<!DOCTYPE html>
<html>
	<body>
		<ul>
			{{ range $key, $value := .Packages }}
				<li>{{ $key }} - {{ $value.Repo }}</li>
			{{ end }}
		</ul>
	</body>
</html>
