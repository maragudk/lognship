.PHONY: lint

lint:
	find lexicons -name '*.json' -print0 | xargs -0 goat lex parse
	find lexicons -name '*.json' -print0 | xargs -0 goat lex lint
