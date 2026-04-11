import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/httpx.dart' as httpx;
import './api.dart' as api;

class ManualConfiguration extends StatefulWidget {
  final void Function()? retry;
  final void Function(api.Daemon) connect;
  final Alignment? alignment;

  ManualConfiguration({
    super.key,
    this.retry,
    required this.connect,
    this.alignment = Alignment.topCenter,
  });

  @override
  State<ManualConfiguration> createState() => _ManualConfigurationView();
}

class _ManualConfigurationView extends State<ManualConfiguration> {
  final String defaultLocalhost = httpx.localhost();
  Widget _cause = ds.Error.zero;
  String _hostname = '';

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return forms.Container(
      alignment: widget.alignment,
      padding: defaults.padding,
      ds.ErrorScreen(
        cause: _cause,
        Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            forms.Field(
              label: SelectableText("hostname"),
              input: TextFormField(
                autofocus: true,
                decoration: new InputDecoration(
                  hintText: defaultLocalhost,
                  helperText: "hostname and port for the retrovibed instance",
                ),
                onChanged: (v) {
                  setState(() {
                    _hostname = v;
                  });
                },
              ),
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                if (widget.retry != null)
                  ds.LoadingButton(
                    Text("retry"),
                    onPressed: () async => widget.retry!(),
                  ),
                ds.LoadingButton(
                  Text("connect"),
                  onPressed: () {
                    return api.daemons
                        .create(
                          api.DaemonCreateRequest(
                            daemon: api.Daemon(
                              hostname:
                                  _hostname.isEmpty
                                      ? defaultLocalhost
                                      : _hostname,
                            ),
                          ),
                        )
                        .then((d) {
                          return widget.connect(d.daemon);
                        })
                        .catchError((cause) {
                          setState(() {
                            _cause = ds.Errors.httpauto(
                              cause,
                              onTap: _reseterr,
                            );
                          });
                        }, test: httpx.ErrorsTest.httpauto)
                        .catchError((cause) {
                          setState(() {
                            _cause = ds.Error.unknown(cause, onTap: _reseterr);
                          });
                        });
                  },
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
