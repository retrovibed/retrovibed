import 'dart:async';
import 'dart:collection';

import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/ddisc.dart' as ddisc;
import 'discovery.card.dart';
import 'empty.results.dart';
import 'search.button.dart';
import 'package:retrovibed/design.kit/stateful.dart';

class DiscoveryGrid extends StatefulWidget {
  final ValueNotifier<media.MediaSearchState> search;
  final List<Widget> leading;

  const DiscoveryGrid({super.key, required this.search, this.leading = const []});

  @override
  State<DiscoveryGrid> createState() => _DiscoveryGridState();
}

class _DiscoveryGridState extends State<DiscoveryGrid> with LoadingState {
  Widget _cause = ds.Error.zero;
  bool _loading = false;
  List<ddisc.Discovery> _items = [];
  media.MediaSearchRequest? _lastFetchedNext;
  StreamSubscription<ddisc.Discovery>? _subscription;
  final SplayTreeSet<ddisc.Discovery> _found = SplayTreeSet<ddisc.Discovery>(ddisc.compare);

  void reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  Future<void> refresh({bool triggered = true}) {
    _subscription?.cancel();
    _found.clear();

    final req = widget.search.value.next;
    final category = mimex.category(req.mimetypes);

    if (category.isEmpty) {
      setState(() {
        _items = [];
      });
      widget.search.value = media.MediaSearchState(next: req, count: 0);
      return Future.value();
    }

    setState(() {
      _items = [];
      _cause = ds.Error.zero;
      _loading = widget.search.value.next.query.isNotEmpty;
    });

    return httpx
        .withRetry(
          () => ddisc.api.locate(
            media.Locate(query: req.query, mimetype: category, adult: req.adult),
            options: [authn.request(authn.AuthzCache.meta(context))],
          ),
        )
        .then((stream) {
          final done = Completer<void>();
          _subscription = stream.listen(
            (item) {
              _found.add(item);
              setState(() {
                _items = _found.toList();
                _loading = false;
              });
              widget.search.value = media.MediaSearchState(next: req, count: _found.length);
            },
            cancelOnError: true,
            onError: done.completeError,
            onDone: done.complete,
          );
          return done.future;
        })
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unknown(cause, onTap: reseterr);
          });
        })
        .whenComplete(() {
          setState(() {
            _loading = false;
          });
        });
  }

  @override
  void dispose() {
    _subscription?.cancel();
    super.dispose();
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

        return ds.ErrorScreen(
          cause: _cause,
          ds.Grid<ddisc.Discovery>(
            key: ValueKey('discovery.grid'),
            (context, v) => DiscoveredCard(v),
            children: _items,
            loading: _loading,
            physics: AlwaysScrollableScrollPhysics(),
            leading: widget.leading,
            empty: EmptyResults(SearchButton(search: state)),
          ),
        );
      },
    );
  }
}
