import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'known.media.card.dart';
import './api.dart' as api;

class KnownMediaLocator extends StatefulWidget {
  final api.Known current;
  final Future<api.LocateCreateResponse> Function(api.Locate req, {List<httpx.Option> options}) locate;
  final IconData icon;
  final Widget help;
  final Widget? trailing;

  const KnownMediaLocator(
    this.current, {
    super.key,
    this.locate = api.locate.create,
    this.icon = Icons.search,
    this.help = ds.HelpScope.None,
    this.trailing,
  });

  @override
  State<StatefulWidget> createState() => _KnownMediaLocator();
}

class _KnownMediaLocator extends State<KnownMediaLocator> {
  bool _loading = false;
  Widget _cause = ds.Error.zero;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  void _onTap() async {
    setState(() {
      _loading = true;
      _cause = ds.Error.zero;
    });

    widget
        .locate(
          api.Locate.create()..knownMediaId = widget.current.id,
          options: [authn.request(authn.AuthzCache.meta(context))],
        )
        .then((_) {
          setState(() {
            _loading = false;
          });
        })
        .catchError((e) {
          setState(() {
            _loading = false;
            _cause = ds.Error.unknown(e, onTap: reseterr);
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    return ds.Loading(
      loading: _loading,
      cause: _cause,
      KnownMediaCard(
        widget.current,
        icon: widget.icon,
        help: widget.help,
        onTap: _loading ? null : _onTap,
        trailing: widget.trailing,
      ),
    );
  }
}
