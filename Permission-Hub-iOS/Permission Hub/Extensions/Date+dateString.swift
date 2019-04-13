//
//  Date+dateString.swift
//  Permission Hub
//
//  Created by Corné on 13/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import Foundation

extension Date {

    func dateAndTimeString(ofStyle style: DateFormatter.Style = .short) -> String {
        return "\(dateString(ofStyle: style)) / \(timeString(ofStyle: style))"
    }

    func dateString(ofStyle style: DateFormatter.Style = .short) -> String {

        if Calendar.current.isDateInToday(self) {
            return "Today"
        } else if Calendar.current.isDateInYesterday(self) {
            return "Yesterday"
        } else {
            return DateFormatter.localizedString(
                from: self,
                dateStyle: style,
                timeStyle: .none)
        }
    }

    func timeString(ofStyle style: DateFormatter.Style = .short) -> String {

        return DateFormatter.localizedString(
            from: self,
            dateStyle: .none,
            timeStyle: style)
    }
}
