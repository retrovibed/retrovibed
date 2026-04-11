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
    return Column(
      mainAxisSize: MainAxisSize.max,
      crossAxisAlignment: CrossAxisAlignment.stretch,
      spacing: defaults.spacing,
      children: [
        ds.Card(Current()),
        ds.Card(AuthzMetaDisplay.current()),
        ds.Card(const AuthzDeeppool()),
      ],
    );
  }
}
