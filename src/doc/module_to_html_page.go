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

			body {
				max-width: 40em;
			}
			
			code {
				font-family: 'Source Code Pro', 'Ubuntu Mono', 'Liberation Mono', Courier, monospace;
			}
			
			div.description p strong {
				color: #a3a37a;
			}
			
			div.description p code {
				color: #939393;
			}

			div.prototype {
				padding-left: 1.5em;
				text-indent: -1.5em;
			}
			
			code.params {
				color: #9f9f9f;
				font-size: small;	
			}
			
			.paren {
				color: #767676;
			}

			.comma {
				color: #c1c1c1;
			}

			.alias {
			  color: #cd99cd;
			}
			
			.arrow {
				color: #979797;
			}
			
			.primitive {
				color: #548554;
			}
			
			.customtype {
				color: #6babab;
			}

			.localtype {
				color: #0893b0;
			}
			
			.invokertype {
				color: #a6abff;
			}

			.functiontype {
				color: #a6ab5f;
			}

			.unmanagedtype {
				color: #ab4e6a;
			}

			.recordtype {
				color: #ab9b75;
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
