import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './api.dart' as api;
import './daemon.auto.dart';

class DaemonDropdownItem extends StatefulWidget {
  final api.Daemon daemon;
  final VoidCallback onTap;

  const DaemonDropdownItem({
    super.key,
    required this.daemon,
    required this.onTap,
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

    return ds.ErrorScreen(
      cause: _cause,
      tint: defaults.dangerTint,
      InkWell(
        onTap: () {
          EndpointAuto.of(
            context,
          )?.refreshNoErrHandling(Future.value(widget.daemon)).then((_) => widget.onTap()).catchError((e) {
            setState(() {
              _cause = ds.Error.unknown(e, onTap: reseterr);
            });
          });
        },
        child: Container(
          padding: defaults.padding,
          child: Text(widget.daemon.description),
        ),
      ),
    );
  }
}
