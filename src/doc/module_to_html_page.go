package doc

import (
	"fmt"
	"io"

	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/compiler/src/loader"
)

func PackagesToHtmlPage(writer io.Writer, packages []*loader.Package) error {
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
				max-width: 40rem;
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

			div.description p {
				line-height: 1.3;
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

			.keyword {
				color: #f2c9c9;
			}

			.alias {
			  color: #cd99cd;
			}
			
			.arrow {
				color: #979797;
			}
			
			.primitivetype {
				color: #9dd0b3;
			}
			
			.customtype {
				color: #6babab;
			}

			.customtypename {
				color: #6babab;
			}

			.customtypevariant {
				color: #6bfefe;
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

			.unmanagedname {
				color: #ab2e2a;
			}

			.recordtype {
				color: #ab9b75;
			}

			.recordtypefield {
				color: #919191;
			}

			.operator {
				color: #eb8ae5;
			}

			.keyword {
				color: #ffecec;
			}

			.typesymbol {
				color: #d2c19a;
			}

			.number {
				color: #efe48b;
			}

			.recordfield {
				color: #919191;
			}

			.modulereference {
				color: #a37df6;
			}

			.typegenerator {
				color: #5d96b8;
			}
			
			.swamp, .swamp-function-prototype, .swamp-value {
				background-color: #404040;
				padding: 0.5rem;
				border-radius: .3rem;
				display: block;
			}

			.admonition {
				border-left: .2rem solid #448aff;
				background-color: red;
				padding: 0.6rem;
				border-radius: 0.3rem;
				margin: 1.5em 0;
			}

			.admonition-title {
				background-color: rgba(244, 222, 109, 0.33);
				margin: 0 -.6rem 0 -.8rem;
				padding: .4rem .6rem .4rem 2rem;
			}

			.warning {
 				background-color:rgb(208, 120, 4);
 				border-color:#ffb554;
			}

			.admonition-title > .warning  {
 				background-color:rgba(255,145,0,.1);
 				border-color:#ff9100;
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

			code a:link {
				text-decoration: none;
			}

			code a:hover {
				text-decoration: none;
			}

			code a:active {
				text-decoration: none;
			}

			code a:visited {
				text-decoration: none;
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

	documentedModules := make(map[string]*decorated.Module)
	for _, compiledPackage := range packages {
		for _, module := range compiledPackage.AllModules() {
			if _, alreadyDocumented := documentedModules[module.FullyQualifiedModuleName().String()]; alreadyDocumented {
				continue
			}
			documentedModules[module.FullyQualifiedModuleName().String()] = module
		}
	}

	keys := make([]string, 0, len(documentedModules))
	for k := range documentedModules {
		keys = append(keys, k)
	}

	for _, key := range keys {
		compiledModule := documentedModules[key]
		if err := ModuleToHtml(writer, compiledModule); err != nil {
			return err
		}
	}

	fmt.Fprintf(writer, footer)

	return nil
}
