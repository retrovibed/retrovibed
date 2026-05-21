import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/authn.dart' as authn;
import 'known.media.card.dart';
import 'known.media.typography.dart';
import './api.dart' as api;

class KnownMediaDropdown extends StatefulWidget {
  final api.FnKnownSearch search;
  final TextEditingController? controller;
  final FocusNode? focus;
  final String current;
  final void Function(api.Known? k)? onChange;
  const KnownMediaDropdown({
    super.key,
    this.search = api.known.search,
    this.controller,
    this.focus,
    this.current = "",
    this.onChange,
  });

  @override
  State<StatefulWidget> createState() => _KnownMediaDropdown();
}

class _KnownMediaDropdown extends State<KnownMediaDropdown> {
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  api.KnownSearchResponse _res = api.known.response(
    next: api.known.request(limit: 4),
  );
  api.Known? current = null;
  List<httpx.Option> _authOptions = [];

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    _authOptions = [authn.request(authn.AuthzCache.meta(context))];
  }

  void reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  Future<void> refresh(api.KnownSearchRequest req) {
    return widget
        .search(req, options: _authOptions)
        .then((v) {
          setState(() {
            _res = v;
            _loading = false;
          });

          widget.focus?.requestFocus();
          ds.textediting.refocus(widget.controller);
        })
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unauthorized(cause, onTap: reseterr);
            _loading = false;
          });
        }, test: httpx.ErrorsTest.unauthorized)
        .catchError((e) {
          setState(() {
            _cause = ds.Error.unknown(e, onTap: reseterr);
            _loading = false;
          });
        });
  }

  @override
  void initState() {
    super.initState();

    if (uuidx.isMinMax(uuidx.fromString(widget.current))) {
      refresh(_res.next);
      return;
    }

    WidgetsBinding.instance.addPostFrameCallback((_) {
      _authOptions = [authn.request(authn.AuthzCache.meta(context))];
      api.known
          .cached(
            widget.current,
            () => api.known.get(widget.current, options: _authOptions),
          )
          .then(
            (w) => setState(() {
              current = w.known;
            }),
          )
          .whenComplete(() => refresh(_res.next));
    });
  }

  @override
  void deactivate() {
    if (current == null) {
      widget.onChange?.call(current);
    }
    super.deactivate();
  }

  @override
  Widget build(BuildContext context) {
    if (current != null) {
      return KnownMediaTypography(
        current!,
        trailing: [
          Spacer(),
          KnownMediaTypography.removebtn(
            context,
            widget.current,
            onPressed:
                () => setState(() {
                  current = null;
                }),
          ),
        ],
      );
    }

    final defaults = ds.Defaults.of(context);

    return Column(
      spacing: defaults.spacing,
      children: [
        ds.Container(
          padding: defaults.padding,
          ds.SearchTray(
            decoration: InputDecoration(hintText: "search known media"),
            controller: widget.controller,
            focus: widget.focus,
            onSubmitted: (v) {
              setState(() {
                _res.next.query = v;
                _res.next.offset = ds.Grid.int64(0);
              });
              return refresh(_res.next);
            },
            next: (i) {
              setState(() {
                _res.next.offset = i;
              });
              refresh(_res.next);
            },
            current: _res.next.offset,
            empty: ds.Grid.int64(_res.items.length) < _res.next.limit,
            autofocus: defaults.desktop,
          ),
        ),
        ds.Loading(
          loading: _res.items.isEmpty,
          overlay: ds.Empty,
          forms.Container(
            ds.Grid(
              children: _res.items,
              loading: _loading,
              cause: _cause,
              leading: [],
              (context, v) => KnownMediaCard(
                v,
                icon: Icons.search,
                onDoubleTap:
                    widget.onChange == null
                        ? null
                        : () {
                          setState(() {
                            current = v;
                          });
                          widget.onChange!(v);
                        },
              ),
            ),
          ),
        ),
      ],
    );
  }
}
