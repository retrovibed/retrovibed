import 'package:flutter/material.dart';
import 'help.dart';

abstract class buttons {
  static Widget search({required VoidCallback? onPressed, double? size, Widget help = HelpScope.None}) {
    return Help(
      IconButton(onPressed: onPressed, icon: Icon(Icons.search_rounded, size: size)),
      help,
    );
  }

  static Widget refresh({required VoidCallback? onPressed, double? size, Widget help = HelpScope.None}) {
    return Help(IconButton(onPressed: onPressed, icon: Icon(Icons.refresh, size: size)), help);
  }

  static Widget settings({required VoidCallback? onPressed, double? size, Widget help = HelpScope.None}) {
    return Help(IconButton(onPressed: onPressed, icon: Icon(Icons.tune, size: size)), help);
  }

  static Widget link({required VoidCallback? onPressed, double? size, Widget help = HelpScope.None}) {
    return Help(IconButton(onPressed: onPressed, icon: Icon(Icons.link, size: size)), help);
  }

  static Widget remove({required VoidCallback? onPressed, double? size, Widget help = HelpScope.None}) {
    return Help(IconButton(onPressed: onPressed, icon: Icon(Icons.clear, size: size)), help);
  }

  static Widget accept({required VoidCallback? onPressed, double? size, Widget help = HelpScope.None}) {
    return Help(IconButton(onPressed: onPressed, icon: Icon(Icons.check, size: size)), help);
  }

  static Widget copy({required VoidCallback? onPressed, double? size, Widget help = HelpScope.None}) {
    return Help(IconButton(onPressed: onPressed, icon: Icon(Icons.copy, size: size)), help);
  }

  static Widget help({required VoidCallback? onPressed, double? size, Widget help = HelpScope.None}) {
    return Help(IconButton(onPressed: onPressed, icon: Icon(Icons.help_center_outlined, size: size)), help);
  }
}
