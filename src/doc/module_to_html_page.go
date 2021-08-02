package doc

import (
	"fmt"
	"io"

	"github.com/swamp/compiler/src/loader"
)

func PackagesToHtmlPage(writer io.Writer, packages []*loader.Package) {
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

			h1 {
				color: #969671;
				margin-top: 1.1em;
			}

			h2 {
				color: #70769a;
				margin-top: 1.1em;
				margin-bottom: 0.5em;
			}

			h3 {
				color: #fafafa;
				margin-top: 2em;
				margin-bottom: 0.5em;
			}
		</style>
	</head>
	<body>
`

	footer := `
	</body>
</html>
`

	fmt.Fprintf(writer, header)

	for _, compiledPackage := range packages {
		fmt.Fprintf(writer, "\n\n\n<hr/>\n<h1>Package %v</h1>\n", compiledPackage.Name())
		for _, module := range compiledPackage.AllModules() {
			ModuleToHtml(writer, module)
		}
	}

	fmt.Fprintf(writer, footer)
}
