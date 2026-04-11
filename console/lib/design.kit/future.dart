import 'package:flutter/material.dart';
import './screens.dart' as screens;

FutureBuilder<T> future<T>(
  T defaults,
  Future<T> pending,
  Widget Function(AsyncSnapshot<T>) render,
) {
  return FutureBuilder(
    initialData: defaults,
    future: pending,
    builder: (BuildContext ctx, AsyncSnapshot<T> snapshot) {
      final effective = AsyncSnapshot<T>.withData(
        snapshot.connectionState,
        snapshot.data ?? defaults,
      );
      return screens.Loading(
        loading: snapshot.connectionState != ConnectionState.done,
        render(effective),
      );
    },
  );
}
