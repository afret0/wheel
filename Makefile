tag = v1.1.882

prod:
	git commit -am "f" && git push
	git tag $(tag)
	git push origin $(tag)


.PHONY: prod