import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/meta/api.dart' as meta;
import 'package:retrovibed/retrovibed.dart' as retro;
import 'package:url_launcher/url_launcher.dart';
import './overview.dart';

class Card extends StatefulWidget {
  final EdgeInsets margin;
  final Function(Widget w)? onPressed;
  final Future<authn.Session> Function(BuildContext) session;
  final Future<meta.Authn> Function() current;
  const Card({
    super.key,
    this.margin = EdgeInsets.zero,
    this.onPressed,
    this.session = authn.Authenticated.session,
    this.current = meta.current,
  });

  @override
  State<Card> createState() => _CardState();
}

class _CardState extends State<Card> {
  String _fallbackUsername = retro.username();
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  authn.Session _session = authn.Session();
  meta.Authn _authn = meta.Authn();

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  @override
  void initState() {
    super.initState();
    _fetch();
  }

  void _fetch() {
    setState(() {
      _loading = true;
      _cause = ds.Error.zero;
    });

    Future.wait([
          httpx.withRetry(() => widget.session(context)),
          httpx.withRetry(() => widget.current()),
        ])
        .then((results) {
          setState(() {
            _loading = false;
            _session = results[0] as authn.Session;
            _authn = results[1] as meta.Authn;
          });
        })
        .catchError((e) {
          setState(() {
            _loading = false;
            _cause = ds.Error.unknown(e, onTap: _fetch);
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);
    final _openWebConsole = () {
      authn
          .otp(options: [authn.DeeppoolAuthzCache.bearer(context)])
          .then((r) {
            final Uri q = Uri.https(httpx.consoleendpoint(), "/", {
              "lt": r.token,
            });
            launchUrl(q);
          })
          .catchError((e) {
            setState(() {
              _cause = ds.Error.unknown(e, onTap: _fetch);
            });
          });
    };

    final tap = widget.onPressed == null ? null : () => widget.onPressed!(const Overview());
    return ds.Card(
      alignment: Alignment.center,
      margin: widget.margin,
      onTap: tap,
      help: ds.Hint(
        label: const Text("Profiles"),
        description: const Text("manage user profiles"),
      ),
      ds.Loading(
        loading: _loading,
        cause: _cause,
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          spacing: defaults.spacing / 4,
          children: [
            Text("Account", style: theme.textTheme.titleMedium),
            ds.Copyable(
              Text(
                [
                  _session.account.description,
                  _session.account.id,
                ].firstWhere((v) => v.isNotEmpty, orElse: () => ''),
                style: theme.textTheme.bodyLarge,
                overflow: TextOverflow.ellipsis,
              ),
              onPressed: ds.Copyable.copy(_session.account.id),
            ),
            ds.Copyable(
              Text(
                [
                  _authn.profile.display,
                  _session.profile.display,
                  _fallbackUsername,
                ].firstWhere((v) => v.isNotEmpty, orElse: () => ''),
                style: theme.textTheme.bodyMedium,
                overflow: TextOverflow.ellipsis,
              ),
              onPressed: ds.Copyable.copy(retro.public_key()),
            ),
            SizedBox(
              width: double.infinity,
              child: OutlinedButton(
                onPressed: _openWebConsole,
                child: Text("Web Console"),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
