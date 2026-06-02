import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'api.dart' as api;
import 'community.pb.dart';
import 'content.detail.dart';

typedef FnPublished =
    Future<PublishedContentListResponse> Function(
      String id, {
      List<httpx.Option> options,
      PublishedContentListRequest req,
    });

typedef FnPublishedTombstone =
    Future<PublishContentDeleteResponse> Function(
      String id, {
      List<httpx.Option> options,
    });

class CommunityContentDisplay extends StatefulWidget {
  final Community community;
  final FnPublished apipublished;
  final FnPublishedTombstone apitombstone;
  final Widget help;

  const CommunityContentDisplay({
    super.key,
    required this.community,
    this.apipublished = api.API.published,
    this.apitombstone = api.API.publishedtombstone,
    this.help = ds.HelpScope.None,
  });

  @override
  State<CommunityContentDisplay> createState() => _CommunityContentDisplayState();
}

class _CommunityContentDisplayState extends State<CommunityContentDisplay> {
  PublishedContentListResponse _resp = PublishedContentListResponse(
    next: PublishedContentListRequest(
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

  Future<void> _refresh(PublishedContentListRequest req) {
    setState(() {
      _loading = true;
      _cause = ds.Error.zero;
    });

    return httpx
        .withRetry(
          () => widget.apipublished(
            widget.community.id,
            req: req,
            options: [authn.request(authn.AuthzCache.meta(context))],
          ),
        )
        .then((response) {
          setState(() {
            _resp = response;
          });
        })
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
    final session = authn.Authenticated.syncSession(context);
    final owned = session.account.id == widget.community.accountId;

    return ds.Table<PublishedContent>(
      ds.Table.expanded<PublishedContent>(
        (item) => _ContentRow(
          community: widget.community,
          item: item,
          onDelete:
              owned
                  ? (content) => httpx
                      .withRetry(
                        () => widget.apitombstone(content.id, options: [authn.request(authn.AuthzCache.meta(context))]),
                      )
                      .then((_) {
                        setState(() {
                          _resp.items.remove(content);
                        });
                      })
                  : null,
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
  final Community community;
  final PublishedContent item;
  final Future<void> Function(PublishedContent)? onDelete;

  const _ContentRow({required this.community, required this.item, this.onDelete});

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
        if (item.archivedId != uuidx.min()) Icon(Icons.archive, color: theme.colorScheme.primary),
        if (onDelete != null)
          ds.LoadingIconButton(
            icon: Icon(Icons.delete, color: Colors.red),
            tooltip: 'remove content',
            onPressed: () async {
              ds.modals
                  .of(context)
                  ?.push(
                    ds.Confirmation.yesNo(
                      content: Text(
                        'Are you sure you want to delist "${item.title.isNotEmpty ? item.title : item.id}"?',
                      ),
                      onConfirm: (context) {
                        onDelete!(item).whenComplete(() {
                          ds.modals.of(context)?.push(null);
                        });
                      },
                      onCancel: (context) {
                        ds.modals.of(context)?.push(null);
                      },
                    ),
                  );
            },
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
