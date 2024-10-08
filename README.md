
# Feedback ingestion

The feedback ingestion service handles pull and push-based subscriptions to ingest feedback from the default sources
### Requirements
- Push and Pull Integration Model ✅
- Metadata Ingestion: each source can have different types of metadata values, e.g. app-version from Playstore, Country from Twitter, etc. ✅
- Multi-tenancy ✅
- Transformation to a uniform internal structure, which should ✅
    - Support different types of feedback data e.g Reviews, Conversations, etc
    - Support source-specific metadata
    - Have common record-level attributes like record language, tenant info, source info, etc

**Good to have**
- Idempotency: Ability to de-dupe ingested feedbacks ✅ - handled as for the same tenantId, source, sub_source the feedback cannot be ingested twice

- Supporting multiple feedback sources of the same type for a tenant, e.g feedback
from two different Playstore Apps for the same tenant ✅ - 2 different subscriptions can be created for the same source and tenant by different sub_source_id

### DB design

![db design](https://github.com/harish-dalal/feedback-ingestion/blob/main/assets/fb_ingest_db.png?raw=true)
### High-level design

![High-level design](https://github.com/harish-dalal/feedback-ingestion/blob/main/assets/fb-ingest-excali.png?raw=true)
### Add source
- Define source and its type in ```pkg/models/source.go```
- Create the new source strategy in the ```feedback-ingestion-system/pkg/integrations ```
- The new integration should implement Push(), Pull(), GetSourceName() and GetSourceType() methods defined in ```pkg/integrations/integration_manager.go```
- Add the new strategy (integration) in the strategiesMap defined in ```pkg/routes/routes.go```
- If the source supports Push (webhook), define the route in ```pkg/routes/routes.go``` for e.g. 
```
// webhooks - need to set webhook routes for all the sources that can support push-based ingestion
	srv.Router.HandleFunc("/webhook/intercom", func(w http.ResponseWriter, r *http.Request) {
		integrationManager.HandleWebhook(w, r, models.SourceIntercom)
	})
```

- And voila !!! we have a new source
## Run Locally

**Prerequisite**
- Install go - https://go.dev/doc/install 
- Install and setup Postgres - https://www.postgresql.org/download/
    - username - **local** and password - **local**
    - create ```fb_ingest``` database
    - run Postgres on _localhost:5432_ (default)


Clone the project

```bash
  git clone https://github.com/harish-dalal/feedback-ingestion.git
```

Go to the project directory

```bash
  cd feedback-ingestion
```

Start the server

```bash
  go run cmd/server/main.go
```

By default, a tenant is created for creating a subscription on it, but a new tenant can also be made by calling the create tenant API provided in the Postman import file below

postman import file - https://drive.google.com/uc?export=download&id=1a_sXBfVT0nU1XuhIL6GXjG0t9IZSrZrI

Default tenant ```cb4d81c7-e1bf-4ca5-900f-665a0e3fc932```

- From the postman APIs 
- ```Call the Subscription/create subscription```
- this will create a pull-based subscription on the default tenant for the source Discourse   

As while starting the server the first cron job would have already been completed, there are 2 ways to see the subscription effect
- Either to wait for 8 hours, so that the next cron job could run. (not possible)
- Restart the service
```bash
  go run cmd/server/main.go
``` 
this will again trigger the first cron job

Check the feedback API - Get feedback for tenant - it should contain 3 feedbacks from discourse (restricting to 3 posts for not overloading the system)

sample response (trimmed to fit)

```json
{
    "data": [
        {
            "id": "706940",
            "tenant_id": "cb4d81c7-e1bf-4ca5-900f-665a0e3fc932",
            "source": "",
            "sub_source_id": "",
            "source_type": "",
            "created_at": "2024-09-04T22:08:51.09762+05:30",
            "updated_at": "2024-09-04T22:08:51.109541+05:30",
            "metadata": null,
            "content": {
                "body": "\u003cp\u003eThe Multilingual Plugin makes ...."
            }
        },
        {
            "id": "695312",
            "tenant_id": "cb4d81c7-e1bf-4ca5-900f-665a0e3fc932",
            "source": "",
            "sub_source_id": "",
            "source_type": "",
            "created_at": "2024-09-04T22:08:51.081291+05:30",
            "updated_at": "2024-09-04T22:08:51.106146+05:30",
            "metadata": null,
            "content": {
                "body": "\u003cp\u003eThe Multilingual Plugin makes ...."
            }
        },
        {
            "id": "603469",
            "tenant_id": "cb4d81c7-e1bf-4ca5-900f-665a0e3fc932",
            "source": "",
            "sub_source_id": "",
            "source_type": "",
            "created_at": "2024-09-04T22:08:51.077169+05:30",
            "updated_at": "2024-09-04T22:08:51.097691+05:30",
            "metadata": null,
            "content": {
                "body": "\u003cp\u003eThe Multilingual Plugin makes ...."
            }
        }
    ]
}
```



## Future scope
- Make cron stateful to handle service restarts and run on multiple separate instances.
- Extract source (source and source type) to be fetched from config or a separate db table to make it less painful to update an existing Source name - which relates to other tables like Subscription.
- more to come...
