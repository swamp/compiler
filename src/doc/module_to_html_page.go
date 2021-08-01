package doc

import (
	"github.com/swamp/compiler/src/loader"
)

func PackagesToHtmlPage(packages []*loader.Package) string {
	header := `
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width,initial-scale=1">
		<title>Swamp Documentation</title>
		<style>
			html {
				color: #eeeeee;
				background: #262626;
				font-family: 'Roboto', 'Open Sans Pro', 'Source Sans Pro', 'Ubuntu', sans-serif;
			}
			
			code {
				font-family: 'Source Code Pro', 'Ubuntu Mono', 'Liberation Mono', Courier, monospace;
			}
			
			div.description p strong {
				color: #f4f4ba;
			}
			
			div.description p code {
				color: #939393;
			}
			
			code.params {
				color: #9f9f9f;
				font-size: small;	
			}
			
			.alias {
			  color: #fa8efa;
			}
			
			.arrow {
				color: #979797;
			}
			
			.primitive {
				color: green;
			}
			
			.customtype {
				color: #44e1e1;
			}
			
			.invoker {
				color: #a6abff;
			}
			
			h3 {
				color: #fafafa;
				margin-top: 2em;
			}
		</style>
	</head>
	<body>
`

	footer := `
	</body>
</html>
`

	segments := ""

	for _, compiledPackage := range packages {
		for _, module := range compiledPackage.AllModules() {
			html := ModuleToHtml(module)
			if len(html) != 0 {
				segments += html
			}
		}
	}

	return header + segments + footer
}
