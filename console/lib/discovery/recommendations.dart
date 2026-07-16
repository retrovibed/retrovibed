import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/library.dart' as lib;
import 'package:retrovibed/langcodex.dart' as langcodex;
import '../media/media.known.pb.dart' as known;

class Recommendations extends StatefulWidget {
  const Recommendations(this.mimetype, {super.key, this.latest = lib.recommendations.latest});

  final String mimetype;
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
    ds.postframe(() => _load());
  }

  @override
  void didUpdateWidget(Recommendations oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.mimetype != widget.mimetype) {
      _load();
    }
  }

  Future<void> _load() async {
    setState(() => _loading = true);
    final auth = authn.request(authn.AuthzCache.meta(context));
    return httpx
        .withRetry(
          () => widget.latest(
            lib.recommendations.request(
              mimetype: widget.mimetype,
              language: langcodex.locale().languageCode,
              adult: false,
            ),
            options: [auth],
          ),
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
              help: ds.Hint(const Text("generate a new (random) recommendation")),
              onPressed: () {
                return httpx
                    .withRetry(
                      () => lib.recommendations.random(
                        lib.RecommendationSearchRequest(
                          mimetype: widget.mimetype,
                          language: langcodex.locale().languageCode,
                          adult: false,
                        ),
                        options: [authn.request(authn.AuthzCache.meta(context))],
                      ),
                    )
                    .catchError((cause) {
                      setState(() {
                        _cause = ds.Errors.httpauto(cause, onTap: reseterr);
                      });
                    }, test: httpx.ErrorsTest.httpauto)
                    .then((_) => _load());
              },
            ),
          ],
        ),
        constraints: BoxConstraints.tightForFinite(height: 256),
        background: ds.Repeat(() => lib.KnownMediaCard(lib.Known(), icon: null)),
        empty: Text(
          'Content partnerships pending',
          style: TextStyle(color: Colors.grey),
        ),
        items: _items
            .map(
              (item) => lib.KnownMediaLocator(
                item,
                icon: Icons.download,
                help: lib.KnownMediaDisplay.hintRecommendations,
                onChange: (v) {
                  final upd = ds.fnOnChange(_items, v, (o) => o.id == item.id);
                  setState(() {
                    _items = upd;
                  });
                },
                trailing: Row(
                  children: [
                    Spacer(),
                    ds.LoadingIconButton.info(
                      iconSize: 18.0,
                      onPressed: () => ds.asyncfn<void>(
                        context,
                        (completion) {
                          return ds.Confirmation.info(
                            padding: EdgeInsets.zero,
                            content: LayoutBuilder(
                              builder: (context, constraints) => lib.KnownMediaCard(
                                item,
                                constraints: BoxConstraints(
                                  maxWidth: constraints.maxWidth < 512 ? constraints.maxWidth * 0.8 : 512,
                                ),
                              ),
                            ),
                            done: completion.complete,
                          );
                        },
                      ),
                    ),
                  ],
                ),
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
