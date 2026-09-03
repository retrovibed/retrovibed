import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/ddisc.dart' as ddisc;
import 'package:retrovibed/discovery.dart' as disc;
import 'api.dart' as api;
import 'known.media.card.dart';

class KnownMediaLocator extends StatefulWidget {
  final api.Known current;
  final Future<api.LocateCreateResponse> Function(api.Locate req, {List<httpx.Option> options}) locate;
  final Future<ddisc.DiscoveryDownloadResponse> Function(String id, {List<httpx.Option> options}) download;
  final Future<bool> Function(BuildContext context, {List<httpx.Option> options}) ensureP2P;
  final Future<api.RecommendationDeleteResponse> Function(String id, {List<httpx.Option> options}) delete;
  final void Function(api.Known? v) onChange;
  final IconData icon;
  final Widget help;
  final List<Widget> leading;
  final List<Widget> trailing;

  const KnownMediaLocator(
    this.current, {
    super.key,
    this.onChange = ds.fnNoop,
    this.ensureP2P = disc.ensureP2P,
    this.locate = api.locate.create,
    this.download = ddisc.api.download,
    this.delete = api.recommendations.delete,
    this.icon = Icons.download_rounded,
    this.help = ds.HelpScope.None,
    this.leading = const [],
    this.trailing = const [],
  });

  static Widget future(
    Future<api.Known> pending, {
    Key? key,
    void Function(api.Known? v) onChange = ds.fnNoop,
    Future<bool> Function(BuildContext context, {List<httpx.Option> options}) ensureP2P = disc.ensureP2P,
    Future<api.LocateCreateResponse> Function(api.Locate req, {List<httpx.Option> options}) locate = api.locate.create,
    Future<ddisc.DiscoveryDownloadResponse> Function(String id, {List<httpx.Option> options}) download =
        ddisc.api.download,
    Future<api.RecommendationDeleteResponse> Function(String id, {List<httpx.Option> options}) delete =
        api.recommendations.delete,
    IconData icon = Icons.download_rounded,
    Widget help = ds.HelpScope.None,
    List<Widget> leading = const [],
    List<Widget> trailing = const [],
  }) {
    return FutureBuilder<api.Known>(
      future: pending,
      builder: (context, snapshot) {
        return ds.Loading(
          loading: !(snapshot.hasData || snapshot.hasError),
          cause: snapshot.hasError ? ds.Errors.httpauto(snapshot.error!) : ds.Error.zero,
          snapshot.data == null
              ? const SizedBox()
              : KnownMediaLocator(
                  snapshot.data!,
                  key: key,
                  onChange: onChange,
                  ensureP2P: ensureP2P,
                  locate: locate,
                  download: download,
                  delete: delete,
                  icon: icon,
                  help: help,
                  trailing: trailing,
                  leading: leading,
                ),
        );
      },
    );
  }

  @override
  State<StatefulWidget> createState() => _KnownMediaLocator();
}

class _KnownMediaLocator extends State<KnownMediaLocator> with ds.LoadingState {
  bool _queued = false;

  @override
  void initState() {
    super.initState();
    loading = false;
  }

  void _onTap() async {
    setState(() {
      loading = true;
      cause = ds.Error.zero;
    });

    final options = [authn.request(authn.AuthzCache.meta(context))];

    widget
        .ensureP2P(context, options: options)
        .then((proceed) {
          if (!proceed) {
            setState(() {
              loading = false;
            });
            return null;
          }

          switch (widget.current.source) {
            case ddisc.sources.discovered:
              return httpx.withRetry(
                () => widget
                    .download(
                      widget.current.uid,
                      options: options,
                    )
                    .then((v) {
                      return widget
                          .delete(widget.current.id, options: options)
                          .catchError((e) => api.RecommendationDeleteResponse.create(), test: httpx.ErrorsTest.err404);
                    })
                    .then((v) {
                      widget.onChange(null);
                      setState(() {
                        _queued = true;
                        loading = false;
                      });
                    }),
              );
            default:
              return httpx.withRetry(
                () => widget
                    .locate(
                      api.Locate.create()
                        ..knownMediaId = widget.current.uid
                        ..mimetype = widget.current.mimetype
                        ..adult = widget.current.adult,
                      options: options,
                    )
                    .then((v) {
                      return widget
                          .delete(widget.current.id, options: options)
                          .catchError((e) => api.RecommendationDeleteResponse.create(), test: httpx.ErrorsTest.err404);
                    })
                    .then((v) {
                      widget.onChange(null);
                      setState(() {
                        _queued = true;
                        loading = false;
                      });
                    }),
              );
          }
        })
        .catchError((e) {
          setState(() {
            loading = false;
            cause = ds.Errors.httpauto(e, onTap: reseterr);
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((e) {
          setState(() {
            loading = false;
            cause = ds.Error.unknown(e, onTap: reseterr);
          });
        });
  }

  void _onPress() async {
    setState(() {
      loading = true;
      cause = ds.Error.zero;
    });

    final options = [authn.request(authn.AuthzCache.meta(context))];

    httpx
        .withRetry(
          () => widget.delete(
            widget.current.id,
            options: options,
          ),
        )
        .then((v) {
          widget.onChange(null);
          setState(() {
            loading = false;
          });
        })
        .catchError((e) {
          widget.onChange(null);
          setState(() {
            loading = false;
          });
        }, test: httpx.ErrorsTest.err404)
        .catchError((e) {
          setState(() {
            loading = false;
            cause = ds.Errors.httpauto(e, onTap: reseterr);
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((e) {
          setState(() {
            loading = false;
            cause = ds.Error.unknown(e, onTap: reseterr);
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    return ds.Loading(
      loading: loading,
      cause: cause,
      KnownMediaCard(
        widget.current,
        icon: _queued ? Icons.query_builder_rounded : widget.icon,
        help: widget.help,
        onTap: loading || _queued ? null : _onTap,
        onLongPress: loading ? null : _onPress,
        leading: widget.leading,
        trailing: widget.trailing,
      ),
    );
  }
}
