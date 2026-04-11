import 'package:flutter/material.dart';

abstract class buttons {
  static IconButton search({required VoidCallback? onPressed, double? size}) {
    return IconButton(onPressed: onPressed, icon: Icon(Icons.search_rounded, size: size));
  }

  static IconButton refresh({required VoidCallback? onPressed, double? size}) {
    return IconButton(onPressed: onPressed, icon: Icon(Icons.refresh, size: size));
  }

  static IconButton settings({required VoidCallback? onPressed, double? size}) {
    return IconButton(onPressed: onPressed, icon: Icon(Icons.tune, size: size));
  }

  static IconButton link({required VoidCallback? onPressed, double? size}) {
    return IconButton(onPressed: onPressed, icon: Icon(Icons.link, size: size));
  }

  static IconButton remove({required VoidCallback? onPressed, double? size}) {
    return IconButton(onPressed: onPressed, icon: Icon(Icons.clear, size: size));
  }

  static IconButton accept({required VoidCallback? onPressed, double? size}) {
    return IconButton(onPressed: onPressed, icon: Icon(Icons.check, size: size));
  }

  static IconButton copy({required VoidCallback? onPressed, double? size}) {
    return IconButton(onPressed: onPressed, icon: Icon(Icons.copy, size: size));
  }
}
