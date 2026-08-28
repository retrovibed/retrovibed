import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './current.dart';
import './authz.meta.display.dart';
import './authz.deeppool.dart';

class Overview extends StatelessWidget {
  const Overview({super.key});

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return ds.Container(
      Column(
        mainAxisSize: MainAxisSize.max,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        spacing: defaults.spacing,
        children: [
          Current(),
          AuthzMetaDisplay.current(),
          // const AuthzDeeppool(),
        ],
      ),
    );
  }
}
