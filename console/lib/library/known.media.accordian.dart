import 'dart:async';
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart' as api;
import 'known.media.typography.dart';
import 'known.media.edit.dart';

class KnownMediaAccordian extends StatelessWidget {
  final Future<api.Known> pending;

  const KnownMediaAccordian(
    this.pending, {
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<api.Known>(
      initialData: api.Known(),
      future: pending,
      builder: (context, snapshot) {
        final defaults = ds.Defaults.of(context);
        return ds.Loading(
          loading: !(snapshot.hasData || snapshot.hasError),
          ds.Container(
            ds.Accordion(
              description: KnownMediaTypography(
                snapshot.data ?? api.Known(),
                decoration: BoxDecoration(color: Colors.transparent),
              ),
              content: KnownMediaEdit(
                snapshot.data ?? api.Known(),
                padding: EdgeInsets.symmetric(vertical: defaults.padding.vertical / 2),
              ),
            ),
          ),
        );
      },
    );
  }
}
