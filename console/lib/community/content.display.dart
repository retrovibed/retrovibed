import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/media/api.dart' as media_api;
import 'package:retrovibed/media/media.pb.dart' as media_pb;
import 'api.dart' as api;
import 'community.pb.dart';
import 'publish.container.dart';

class CommunityContentDisplay extends StatefulWidget {
  final Community community;

  const CommunityContentDisplay({super.key, required this.community});

  @override
  State<CommunityContentDisplay> createState() => _CommunityContentDisplayState();
}

class _CommunityContentDisplayState extends State<CommunityContentDisplay> {
  List<PublishedContent> _content = [];
  bool _loading = true;
  Widget _cause = ds.Error.zero;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _clearCause() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  @override
  void initState() {
    super.initState();
    _loadContent();
  }

  Future<Null> _loadContent() {
    setState(() {
      _loading = true;
      _cause = ds.Error.zero;
    });

    // TODO: deeppool endpoints for published content
    final authOptions = [authn.DeeppoolAuthzCache.bearer(context)];

    return httpx
        .withRetry(
          () => api.API.published(
            widget.community.id,
            options: authOptions,
          ),
        )
        .then((response) {
          setState(() {
            _content = response.items;
          });
        })
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unknown(cause, onTap: _clearCause);
          });
        })
        .whenComplete(() {
          setState(() {
            _loading = false;
          });
        });
  }

  void _showPublishModal() {
    ds.modals.asyncfn<void>(
      context,
      (completion) => PublishContainer(
        onCancel: completion.complete,
        onPublished: () => _loadContent().then(completion.complete),
        community: widget.community,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    final emptyState = Center(
      child: Padding(
        padding: EdgeInsets.all(defaults.spacing * 4),
        child: Column(
          children: [
            Icon(Icons.folder_off, size: 48, color: theme.colorScheme.outline),
            SizedBox(height: defaults.spacing),
            Text(
              'No content published yet',
              style: theme.textTheme.bodyMedium?.copyWith(
                color: theme.colorScheme.outline,
              ),
            ),
            SizedBox(height: defaults.spacing),
            TextButton.icon(
              onPressed: _showPublishModal,
              icon: Icon(Icons.add),
              label: Text('Publish your first content'),
            ),
          ],
        ),
      ),
    );

    final contentList = Expanded(
      child: ListView.builder(
        padding: defaults.padding,
        itemCount: _content.length,
        itemBuilder: (context, index) {
          final item = _content[index];
          return _ContentRow(community: widget.community, item: item);
        },
      ),
    );

    return ds.Loading(
      loading: _loading,
      cause: _cause,
      Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Padding(
            padding: defaults.padding,
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text('Published Content', style: theme.textTheme.titleMedium),
                ElevatedButton.icon(
                  onPressed: _showPublishModal,
                  icon: Icon(Icons.add),
                  label: Text('Publish'),
                ),
              ],
            ),
          ),
          if (_content.isEmpty) emptyState else contentList,
        ],
      ),
    );
  }
}

class _ContentRow extends StatelessWidget {
  final Community community;
  final PublishedContent item;

  const _ContentRow({required this.community, required this.item});

  Future<void> _archive(BuildContext context) {
    final req = media_pb.MagnetCreateRequest(uri: item.magnetUri);
    return httpx.withRetry(
      () => media_api.discovered.magnet(req, options: [authn.AuthzCache.bearer(context)]),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    return Card(
      margin: EdgeInsets.only(bottom: defaults.spacing),
      child: Padding(
        padding: defaults.padding,
        child: Row(
          children: [
            Icon(Icons.movie, size: 40, color: theme.colorScheme.primary),
            SizedBox(width: defaults.spacing),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    item.knownMediaId,
                    style: theme.textTheme.bodyLarge,
                    overflow: TextOverflow.ellipsis,
                  ),
                  SizedBox(height: defaults.spacing / 2),
                  Text(
                    item.magnetUri,
                    style: theme.textTheme.bodySmall?.copyWith(
                      color: theme.colorScheme.outline,
                    ),
                    overflow: TextOverflow.ellipsis,
                  ),
                ],
              ),
            ),
            ds.LoadingIconButton(
              tooltip: 'Archive this content',
              onPressed: () => _archive(context),
              icon: Icon(Icons.download),
            ),
          ],
        ),
      ),
    );
  }
}
