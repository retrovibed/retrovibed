import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './api.dart' as api;

Future<void> _noop() => Future.value();

abstract class Icons {
  static ds.LoadingIconButton delete(
    api.Wireguard current, {
    ds.AsyncVoidCallback? onPressed,
    Key? key,
    String? tooltip,
  }) {
    return ds.LoadingIconButton.delete(
      key: key,
      tooltip: tooltip,
      onPressed: onPressed ?? () => api.wireguard.delete(current.id).then((_) {}),
    );
  }

  static ds.LoadingIconButton edit({
    ds.AsyncVoidCallback? onPressed,
    Key? key,
    String? tooltip,
  }) {
    return ds.LoadingIconButton.edit(
      key: key,
      tooltip: tooltip ?? "edit",
      onPressed: onPressed ?? _noop,
    );
  }

  static ds.LoadingIconButton create({
    ds.AsyncVoidCallback? onPressed,
    Key? key,
    String? tooltip,
  }) {
    return ds.LoadingIconButton.create(
      key: key,
      tooltip: tooltip,
      onPressed: onPressed ?? _noop,
    );
  }
}
