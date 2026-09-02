#### Social Publishing

framework for publishing content to social networks built on top of communities.

#### Changes needed

##### Search

Split community search into two endpoints. the long term plan is that communities are represented locally and synced with the archive. We're part way through improving that.

- rename search to discover for the current search archive.
- reimplement search to only hit the local database.

##### We'll need a specialized id for deeppool publishing.

We want to attach builtin native publishing implementations to the community_publisher. That way we can configure multiple flows.
