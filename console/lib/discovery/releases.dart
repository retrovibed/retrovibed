import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/library.dart' as lib;

class NewReleases extends StatefulWidget {
  const NewReleases(String this.mimetype, {super.key, this.latest = lib.known.latest});

  final String mimetype;
  final lib.FnKnownLatest latest;

  @override
  State<NewReleases> createState() => _NewReleasesState();
}

class _NewReleasesState extends State<NewReleases> {
  Widget _cause = ds.Error.zero;
  bool _loading = true;
  lib.KnownLatestResponse _result = lib.KnownLatestResponse();

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
    WidgetsBinding.instance.addPostFrameCallback((_) => _load());
  }

  @override
  void didUpdateWidget(NewReleases oldWidget) {
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
          () => widget.latest(lib.known.latestRequest(mimetype: widget.mimetype), options: [auth]),
        )
        .then(
          (resp) => setState(() {
            _result = resp;
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
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unknown(cause, onTap: reseterr);
            _loading = false;
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    return ds.CarouselRow(
      title: const Text('New Releases'),
      constraints: BoxConstraints.tightForFinite(height: 256),
      background: ds.Repeat(() => lib.KnownMediaCard(lib.Known(), icon: null)),
      empty: Text(
        'Content partnerships in progress',
        style: TextStyle(color: Colors.grey),
      ),
      items:
          _result.items
              .map(
                (item) => lib.KnownMediaCard(
                  item,
                  icon: Icons.download,
                  help: lib.KnownMediaDisplay.hintReleases,
                  onTap: () {
                    ds.modals.asyncfn(
                      context,
                      (completion) => ds.Confirmation.ok(
                        content: Text("automatic media discovery is not yet implemented"),
                        onConfirm: (_) => completion.complete(),
                        onCancel: (_) => completion.complete(),
                      ),
                    );
                  },
                ),
              )
              .toList(),
      loading: _loading,
      cause: _cause,
    );
  }
}
