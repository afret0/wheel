tag = v1.1.88

prod:
	git commit -am "f" && git push
	git tag $(tag)
	git push origin $(tag)


.PHONY: prod