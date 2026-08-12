import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:url_launcher/url_launcher.dart';
import 'package:retrovibed/design.kit/stateful.dart';

class Current extends StatefulWidget {
  @override
  State<StatefulWidget> createState() {
    return _CurrentState();
  }
}

class _CurrentState extends State<Current> with LoadingState {
  bool _loading = true;
  Widget _cause = ds.Error.zero;
  authn.Session current = authn.Session();

  void refresh() {
    setState(() {
      _loading = true;
      _cause = ds.Error.zero;
    });

    authn.Authenticated.session(context)
        .then((session) {
          setState(() {
            _loading = false;
            current = session;
          });
        })
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.offline(cause, onTap: refresh);
          });
        }, test: ds.ErrorTests.offline)
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.connectivity(cause, onTap: refresh);
          });
        }, test: ds.ErrorTests.connectivity)
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unknown(cause, onTap: refresh);
          });
        })
        .whenComplete(() {
          setState(() {
            _loading = false;
          });
        });
  }

  @override
  void initState() {
    super.initState();
    refresh();
  }

  @override
  Widget build(BuildContext context) {
    return forms.Container(
      padding: EdgeInsets.symmetric(horizontal: 10),
      ds.Loading(
        cause: _cause,
        loading: _loading,
        Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            forms.Field(label: Text("id"), input: Text(current.account.id)),
            forms.Field(
              label: Text("name"),
              input: Text(current.account.description),
            ),
            if (authn.developer(context).subscription)
              TextButton(
                child: Text("open web console"),
                onPressed: () {
                  httpx
                      .withRetry(
                        () => authn.otp(options: [authn.DeeppoolAuthzCache.bearer(context)]),
                      )
                      .then((r) {
                        final Uri q = Uri.https(httpx.consoleendpoint(), "/", {
                          "lt": r.token,
                        });
                        launchUrl(q);
                      })
                      .catchError((cause) {
                        setState(() {
                          _cause = ds.Error.unknown(cause, onTap: refresh);
                        });
                      });
                },
              ),
          ],
        ),
      ),
    );
  }
}
