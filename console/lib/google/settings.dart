import 'dart:async';

import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/google/api.dart' as api;
import 'package:url_launcher/url_launcher.dart';

class Settings extends StatefulWidget {
  const Settings({super.key});

  @override
  State<Settings> createState() => _SettingsState();
}

class _SettingsState extends State<Settings> {
  bool _loading = true;
  Timer? _poll;
  Widget _cause = ds.Error.zero;
  api.YouTubeStatus _youtube = api.YouTubeStatus();

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  @override
  void initState() {
    super.initState();
    ds.postframe(() => _fetch());
  }

  @override
  void dispose() {
    _poll?.cancel();
    super.dispose();
  }

  void _fetch() {
    setState(() {
      _loading = true;
      _cause = ds.Error.zero;
    });

    final auth = [authn.request(authn.AuthzCache.meta(context))];
    httpx
        .withRetry(() => api.YouTube.status(options: auth))
        .then((v) {
          if (v.linked) _poll?.cancel();
          setState(() {
            _loading = false;
            _youtube = v;
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

    final auth = [authn.request(authn.AuthzCache.meta(context))];
    return ds.Loading(
      loading: _loading,
      cause: _cause,
      Padding(
        padding: defaults.padding,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          spacing: defaults.spacing,
          children: [
            Text("YouTube", style: theme.textTheme.titleMedium),
            Text(
              _youtube.linked ? "Connected" : "Not connected",
              style: theme.textTheme.bodySmall,
            ),
            SizedBox(
              width: double.infinity,
              child:
                  _youtube.linked
                      ? ds.LoadingButton(
                        Text("Unlink"),
                        onPressed: () {
                          return httpx
                              .withRetry(
                                () => api.YouTube.unlink(options: auth),
                              )
                              .then((_) => _fetch())
                              .catchError((e) {
                                setState(() {
                                  _cause = ds.Error.unknown(e, onTap: _fetch);
                                });
                              });
                        },
                      )
                      : OutlinedButton(
                        onPressed: () {
                          authn.otp(options: [authn.DeeppoolAuthzCache.bearer(context)]).then((session) {
                            launchUrl(
                              api.YouTube.authUri(token: session.token),
                            );
                            _poll?.cancel();
                            _poll = Timer.periodic(
                              const Duration(seconds: 3),
                              (_) => _fetch(),
                            );
                          });
                        },
                        child: Text("Link YouTube"),
                      ),
            ),
          ],
        ),
      ),
    );
  }
}
