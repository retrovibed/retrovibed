import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/library.dart' as lib;
import '../media/media.known.pb.dart' as known;

class Recommendations extends StatefulWidget {
  const Recommendations({super.key, this.latest = lib.recommendations.latest});

  final lib.FnRecommendations latest;

  @override
  State<Recommendations> createState() => _RecommendationsState();
}

class _RecommendationsState extends State<Recommendations> {
  Widget _cause = ds.Error.zero;
  bool _loading = true;
  List<known.Known> _items = [];

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() => _loading = true);
    final auth = authn.AuthzCache.bearer(context);
    return httpx
        .withRetry(
          () => widget.latest(lib.recommendations.request(), options: [auth]),
        )
        .then(
          (resp) => setState(() {
            _items = resp.items.toList();
            _loading = false;
          }),
        )
        .catchError((cause) {
          setState(() {
            _loading = false;
          });
        }, test: httpx.ErrorsTest.httpnotimplemented)
        .catchError((cause) {
          setState(() {
            _cause = ds.Errors.httpauto(cause, onTap: reseterr);
            _loading = false;
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((e) {
          setState(() {
            _cause = ds.Error.unknown(e, onTap: reseterr);
            _loading = false;
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    return ds.KeyPressAware.refresh(
      ds.CarouselRow(
        title: Row(
          children: [
            const Text('Recommendations'),
            Spacer(),
            ds.LoadingIconButton.refresh(
              onPressed: () {
                return httpx
                    .withRetry(() => lib.recommendations.random(options: [authn.AuthzCache.bearer(context)]))
                    .then((_) => _load());
              },
            ),
          ],
        ),
        constraints: BoxConstraints.tightForFinite(height: 256),
        background: ds.Repeat(() => lib.KnownMediaCard(lib.Known(), icon: null)),
        empty: Text(
          'Content partnerships in progress',
          style: TextStyle(color: Colors.grey),
        ),
        items:
            _items
                .map(
                  (item) => lib.KnownMediaCard(
                    item,
                    icon: Icons.download,
                    onTap: () {
                      ds.modals.asyncfn(
                        context,
                        (completion) => ds.Confirmation.ok(
                          content: Text("automatic media discovery is not yet implemented"),
                          onConfirm: completion.complete,
                          onCancel: completion.complete,
                        ),
                      );
                    },
                  ),
                )
                .toList(),
        loading: _loading,
        cause: _cause,
      ),
      onPress: _load,
    );
  }
}
