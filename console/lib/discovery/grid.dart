import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/langcodex.dart' as langcodex;
import 'package:retrovibed/library/api.dart' as api;
import 'package:retrovibed/library/known.media.download.list.dart';
import 'package:retrovibed/library/discovery.locator.dart';

class DiscoveryGrid extends StatefulWidget {
  final ValueNotifier<media.MediaSearchState> search;

  const DiscoveryGrid({super.key, required this.search});

  @override
  State<DiscoveryGrid> createState() => _DiscoveryGridState();
}

class _DiscoveryGridState extends State<DiscoveryGrid> {
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  List<api.Known> _items = [];
  media.MediaSearchRequest? _lastFetchedNext;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  Future<void> refresh() {
    final req = widget.search.value.next;
    final category = mimex.category(req.mimetypes);

    if (category.isEmpty) {
      setState(() {
        _items = [];
        _loading = false;
      });
      widget.search.value = media.MediaSearchState(next: req, count: 0);
      return Future.value();
    }

    return httpx
        .withRetry(
          () => api.known.search(
            api.known.request(
              language: langcodex.locale().languageCode,
              mimetype: category,
              adult: req.adult,
              query: req.query,
              limit: req.limit.toInt(),
            ),
            options: [authn.request(authn.AuthzCache.meta(context))],
          ),
        )
        .then((v) {
          setState(() {
            _items = v.items;
            _loading = false;
          });
          widget.search.value = media.MediaSearchState(next: req, count: v.items.length);
        })
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unknown(cause, onTap: reseterr);
            _loading = false;
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder<media.MediaSearchState>(
      valueListenable: widget.search,
      builder: (context, state, _) {
        if (!identical(state.next, _lastFetchedNext)) {
          _lastFetchedNext = state.next;
          ds.postframe(() => refresh());
        }

        final defaults = ds.Defaults.of(context);
        final category = mimex.category(state.next.mimetypes);

        return ds.Loading(
          loading: _loading,
          cause: _cause,
          KnownMediaDownloadList(
            key: ValueKey('library.download.list'),
            children: _items,
            leading: Column(
              children: [
                DiscoveryLocator(
                  key: ValueKey('library.disc.locator'),
                  query: state.next.query.trim(),
                  mimetype: category,
                ),
              ],
            ),
            empty: Center(
              child: Padding(
                padding: defaults.padding,
                child: Text(
                  "no known media matches the results",
                  style: TextStyle(
                    color: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.6),
                  ),
                ),
              ),
            ),
          ),
        );
      },
    );
  }
}
