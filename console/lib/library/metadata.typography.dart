import 'package:flutter/material.dart';

Widget archived(String uid) {
  switch (uid) {
    case "":
    case "00000000-0000-0000-0000-000000000000":
      return Text("disabled");
    case "ffffffff-ffff-ffff-ffff-ffffffffffff":
      return Text("pending");
    default:
      return Text("archived");
  }
}

Widget sharing(String uid) {
  return toggled(uid, Text("personal"), Text("shared"));
}

Widget toggled(String uid, Widget minmax, Widget o) {
  switch (uid) {
    case "":
    case "00000000-0000-0000-0000-000000000000":
    case "ffffffff-ffff-ffff-ffff-ffffffffffff":
      return minmax;
    default:
      return o;
  }
}