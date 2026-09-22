import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/library.dart' as lib;
import 'package:retrovibed/langcodex.dart' as langcodex;

class NewReleases extends StatefulWidget {
  const NewReleases(String this.mimetype, {super.key, this.latest = lib.known.latest});

  final String mimetype;
  final lib.FnKnownLatest latest;

  @override
  State<NewReleases> createState() => _NewReleasesState();
}

class _NewReleasesState extends State<NewReleases> with ds.LoadingState {
  lib.KnownLatestResponse _result = lib.KnownLatestResponse();

  @override
  void initState() {
    super.initState();
    ds.postframe(() => _load());
  }

  @override
  void didUpdateWidget(NewReleases oldWidget) {
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
          () => widget.latest(
            lib.known.latestRequest(
              language: langcodex.locale().languageCode,
              adult: false,
              mimetype: widget.mimetype,
            ),
            options: [auth],
          ),
        )
        .then(
          (resp) => setState(() {
            _result = resp;
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
        .catchError((cause) {
          setState(() {
            this.cause = ds.Error.unknown(cause, onTap: reseterr);
            loading = false;
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    return ds.CarouselRow(
      title: const Text('New Releases'),
      constraints: BoxConstraints.tightForFinite(height: 256),
      background: ds.Repeat(() => lib.KnownMediaCard(lib.Known(), icon: null)),
      empty: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            'media library is initializing...',
            textAlign: TextAlign.center,
            style: TextStyle(color: Colors.grey),
          ),
        ],
      ),
      items: _result.items
          .map(
            (k) => lib.KnownMediaLocator(
              k,
              icon: Icons.download,
              help: lib.KnownMediaDisplay.hintReleases,
              onChange: (v) {
                setState(() {
                  _result = lib.KnownLatestResponse(
                    next: _result.next,
                    items: ds.fnOnChange(_result.items, v, (o) => o.id == k.id),
                  );
                });
              },
            ),
          )
          .toList(),
      loading: loading,
      cause: cause,
    );
  }
}
