import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'api.dart' as api;
import 'daemon.auto.dart';
import 'daemon.typography.dart';

class DaemonDropdownItem extends StatefulWidget {
  final api.Daemon library;
  final VoidCallback onTap;
  final bool readonly;

  const DaemonDropdownItem({
    super.key,
    required this.library,
    required this.onTap,
    this.readonly = false,
  });

  @override
  State<DaemonDropdownItem> createState() => _DaemonDropdownItemState();
}

class _DaemonDropdownItemState extends State<DaemonDropdownItem> {
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

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final isDevice = api.daemons.isLocalDevice(widget.library);
    return ds.ErrorScreen(
      cause: _cause,
      tint: defaults.dangerTint,
      ds.TableRow.single(
        DaemonTypography(
          widget.library,
          trailing: [
            Visibility(
              visible: !(widget.readonly || isDevice),
              maintainSize: true,
              maintainAnimation: true,
              maintainState: true,
              child: ds.LoadingIconButton.delete(
                onPressed: () async {
                  ds.modals.asyncfn<void>(context, (completion) {
                    return ds.Confirmation.yesNo(
                      content: Text(
                        'Delete ${widget.library.description}?',
                      ),
                      onCancel: (_) => completion.complete(),
                      onConfirm: (_) {
                        httpx.withRetry(
                          () => api.daemons.delete(widget.library.id).then((_) {
                            completion.complete();
                          }),
                        );
                      },
                    );
                  });
                },
                tooltip: "remove library",
              ),
            ),
          ],
        ),
        onTap: () {
          EndpointAuto.of(
            context,
          )?.refreshNoErrHandling(Future.value(widget.library)).then((_) => widget.onTap()).catchError((e) {
            setState(() {
              _cause = ds.Error.unknown(e, onTap: reseterr);
            });
          });
        },
      ),
    );
  }
}
