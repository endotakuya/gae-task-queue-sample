PROJECT_NAME="gae-task-queue-sample"

serve:
	dev_appserver.py app.yaml

deploy:
	find build/go/src -type l -delete
	ls -1 ./vendor | xargs -n 1 -I {} ln -s $(abspath .)/vendor/{} $(abspath .)/build/go/src/{}
	GOPATH=$(abspath build/go) gcloud app deploy --project $(PROJECT_NAME)

clean:
	find build/go/src -type l -delete