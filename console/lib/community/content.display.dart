import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/media/api.dart' as media_api;
import 'package:retrovibed/media/media.pb.dart' as media_pb;
import 'api.dart' as api;
import 'community.pb.dart';
import 'publish.container.dart';

typedef FnPublished = Future<PublishedContentListResponse> Function(
  String id, {
  List<httpx.Option> options,
  int offset,
  int limit,
});

typedef FnMagnet = Future<media_pb.MagnetCreateResponse> Function(
  media_pb.MagnetCreateRequest req, {
  List<httpx.Option> options,
});

class CommunityContentDisplay extends StatefulWidget {
  final Community community;
  final FnPublished apipublished;
  final FnMagnet apimagnet;

  const CommunityContentDisplay({
    super.key,
    required this.community,
    this.apipublished = api.API.published,
    this.apimagnet = media_api.discovered.magnet,
  });

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

    return httpx
        .withRetry(
          () => widget.apipublished(
            widget.community.id,
            options: [authn.DeeppoolAuthzCache.bearer(context)],
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

    return ds.Table<PublishedContent>(
      ds.Table.expanded<PublishedContent>(
        (item) => _ContentRow(
          community: widget.community,
          item: item,
          apimagnet: widget.apimagnet,
        ),
      ),
      loading: _loading,
      cause: _cause,
      children: _content,
      leading: Padding(
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
      empty: Center(
        child: Padding(
          padding: EdgeInsets.all(defaults.spacing * 4),
          child: Column(
            mainAxisSize: MainAxisSize.min,
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
      ),
    );
  }
}

class _ContentRow extends StatelessWidget {
  final Community community;
  final PublishedContent item;
  final FnMagnet apimagnet;

  const _ContentRow({required this.community, required this.item, required this.apimagnet});

  Future<void> _archive(BuildContext context) {
    final req = media_pb.MagnetCreateRequest(uri: item.magnetUri);
    return httpx.withRetry(
      () => apimagnet(req, options: [authn.AuthzCache.bearer(context)]),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return ds.TableRow(
      [
        Icon(Icons.movie, size: 40, color: theme.colorScheme.primary),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                item.knownMediaId,
                style: theme.textTheme.bodyLarge,
                overflow: TextOverflow.ellipsis,
              ),
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
    );
  }
}
