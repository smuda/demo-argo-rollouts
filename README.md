# demo-argo-rollouts

## Prepare

Run `make start`.


## Install with deploy

`make app-deploy`

Open http://cider.apps.127.0.0.1.nip.io/ and show colors.

Change app/helm/values.yaml to use green image and `make app-deploy` again.

It's a quick change.

## Install with rollouts

`make app-rollout`

Show there is no deployment

Change app/helm/values.yaml to use green image and `make app-rollout` again.

The traffic stops at 20% (1 pod). Show the steps in the rollout.

Open http://rollouts.app.127.0.0.1.nip.io/rollouts/app-cider and promote.

Show that `kubectl argo rollouts get rollout cider` exists.

## Install with rollouts and metrics

`make app-rollout-metrics`

Run `watch kubectl argo rollouts get rollout cider`

Show metrics analysis templates.

## Error handling

Change app/helm-metrics/values.yaml to use bad-red image and `make app-rollout-metrics` again.

Notice red results and automatic rollback.

## Traffic Management

`make app-rollout-traffic`

Run `watch kubectl argo rollouts get rollout cider`

Notice that despite having two stable pods, traffic goes by percentage (roughly)
