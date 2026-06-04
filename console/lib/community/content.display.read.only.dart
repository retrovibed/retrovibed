import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/mimex.dart' as mimex;
import 'api.dart' as api;
import 'content.detail.dart';

class ContentDisplayReadOnly extends StatefulWidget {
  final api.Community community;
  final api.FnPublishingSearch apipublished;
  final Widget help;

  const ContentDisplayReadOnly({
    super.key,
    required this.community,
    this.apipublished = api.API.published,
    this.help = ds.HelpScope.None,
  });

  @override
  State<ContentDisplayReadOnly> createState() => _ContentDisplayReadOnlyState();
}

class _ContentDisplayReadOnlyState extends State<ContentDisplayReadOnly> {
  api.PublishedContentSearchResponse _resp = api.PublishedContentSearchResponse(
    next: api.PublishedContentSearchRequest(
      offset: ds.Int64(0),
      limit: ds.Int64(20),
      query: "",
    ),
  );
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
    WidgetsBinding.instance.addPostFrameCallback((_) => _refresh(_resp.next));
  }

  Future<void> _refresh(api.PublishedContentSearchRequest req) {
    setState(() {
      _loading = true;
      _cause = ds.Error.zero;
    });

    return httpx
        .withRetry(
          () => widget.apipublished(
            widget.community.id,
            req: req,
            options: [authn.DeeppoolAuthzCache.bearer(context)],
          ),
        )
        .then((response) {
          setState(() {
            _resp = response;
          });
        })
        .catchError((cause) {
          setState(() {});
        }, test: httpx.ErrorsTest.err404)
        .catchError((cause) {
          setState(() {
            _cause = ds.Errors.httpauto(cause, onTap: _clearCause);
          });
        }, test: httpx.ErrorsTest.httpauto)
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

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    return ds.Table<api.PublishedContent>(
      ds.Table.expanded<api.PublishedContent>(
        (item) => _ContentRow(
          community: widget.community,
          item: item,
        ),
      ),
      loading: _loading,
      cause: _cause,
      children: _resp.items,
      leading: ds.SearchTray(
        decoration: InputDecoration(hintText: "search content"),
        onSubmitted: (q) {
          setState(
            () =>
                _resp.next
                  ..query = q
                  ..offset = ds.Int64(0),
          );
          return _refresh(_resp.next);
        },
        next: (i) {
          setState(() => _resp.next..offset = i);
          _refresh(_resp.next);
        },
        current: _resp.next.offset,
        empty: ds.Int64(_resp.items.length) < _resp.next.limit,
        help: widget.help,
      ),
      empty: Center(
        child: Padding(
          padding: EdgeInsets.all(defaults.spacing * 4),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            spacing: defaults.spacing,
            children: [
              Icon(Icons.folder_off, size: 48, color: theme.colorScheme.outline),
              Text(
                'No content published yet',
                style: theme.textTheme.bodyMedium?.copyWith(
                  color: theme.colorScheme.outline,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _ContentRow extends StatelessWidget {
  final api.Community community;
  final api.PublishedContent item;

  const _ContentRow({
    required this.community,
    required this.item,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    return ds.TableRow(
      key: ValueKey(item.id),
      [
        Icon(mimex.icon(item.mimetype), size: 40, color: theme.colorScheme.primary),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                item.title.isNotEmpty ? item.title : item.id,
                style: theme.textTheme.bodyLarge,
                overflow: TextOverflow.ellipsis,
              ),
              Wrap(
                spacing: defaults.spacing,
                children: [
                  ds.Bytes(item.bytes),
                  ds.Timestamp.iso8601(item.publishedAt, leading: Text('published ')),
                ],
              ),
            ],
          ),
        ),
      ],
      expanded: SizedBox(
        width: double.infinity,
        child: ds.Container(
          padding: defaults.padding,
          PublishedContentDetail(item: item),
        ),
      ),
    );
  }
}
