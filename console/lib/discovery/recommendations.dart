import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/library.dart' as lib;
import 'package:retrovibed/langcodex.dart' as langcodex;
import '../media/media.known.pb.dart' as known;

class Recommendations extends StatefulWidget {
  const Recommendations(
    this.mimetype, {
    super.key,
    this.apilatest = lib.recommendations.latest,
    this.apirefresh = lib.recommendations.refresh,
  });

  final String mimetype;
  final lib.FnRecommendations apilatest;
  final lib.FnRecommendationsRequest apirefresh;

  @override
  State<Recommendations> createState() => _RecommendationsState();
}

class _RecommendationsState extends State<Recommendations> with ds.LoadingState {
  List<known.Known> _items = [];

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
    setState(() => loading = true);
    final auth = authn.request(authn.AuthzCache.meta(context));
    return httpx
        .withRetry(
          () => widget.apilatest(
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
            loading = false;
          }),
        )
        .catchError((cause) {
          setState(() {
            loading = false;
          });
        }, test: httpx.ErrorsTest.notimplemented)
        .catchError((cause) {
          setState(() {
            this.cause = ds.Errors.httpauto(cause, onTap: reseterr);
            loading = false;
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((e) {
          setState(() {
            this.cause = ds.Error.unknown(e, onTap: reseterr);
            loading = false;
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
              help: ds.Hint(const Text("kick off a new recommendation cycle")),
              onPressed: () {
                return httpx
                    .withRetry(
                      () => widget.apirefresh(
                        lib.RecommendationRefreshRequest(
                          mimetype: widget.mimetype,
                          language: langcodex.locale().languageCode,
                          adult: false,
                          limit: ds.Int64(5),
                        ),
                        options: [authn.request(authn.AuthzCache.meta(context))],
                      ),
                    )
                    .then((_) => _load())
                    .catchError((cause) {
                      setState(() {
                        this.cause = ds.Errors.httpauto(cause, onTap: reseterr);
                      });
                    }, test: httpx.ErrorsTest.httpauto);
              },
            ),
          ],
        ),
        constraints: BoxConstraints.tightForFinite(height: 256),
        background: ds.Repeat(() => lib.KnownMediaCard(lib.Known(), icon: null)),
        empty: Text(
          'check back later',
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
                leading: [
                  lib.KnownMediaCard.description(item.description),
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
            )
            .toList(),
        loading: loading,
        cause: cause,
      ),
      onPress: _load,
    );
  }
}
