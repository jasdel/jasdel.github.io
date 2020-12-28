update-theme:
	git submodule update --init --recursive themes/beautifulhugo
	git submodule update --init --recursive themes/even

generate-chroma:
	@#https://github.com/halogenica/beautifulhugo
	hugo gen chromastyles --style=trac > static/css/syntax.css
