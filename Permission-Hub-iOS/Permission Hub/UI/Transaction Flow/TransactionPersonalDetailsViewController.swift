//
//  TransactionPersonalDetailsViewController.swift
//  Permission Hub
//
//  Created by Corné on 13/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

final class TransactionPaymentMethodViewController: PHTableViewController {

    // MARK: - Private properties

    private var transaction: TransactionNotification

    // MARK: - Initialization

    init(transaction: TransactionNotification) {
        self.transaction = transaction

        super.init(title: "", items: [
            .notification(
                type: .verification,
                text: "This company is verified"),
            .description(
                date: transaction.date,
                title: "Payment method",
                description: "Please select out your payment method"),
            .selection(options: [
                "Debit card / Credit card",
                "Paypal",
                "Directly from account"
            ])
        ])
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Life cycle

    override func viewDidLoad() {
        super.viewDidLoad()

        tableView.allowsSelection = false
    }
}

final class TransactionPersonalDetailsViewController: PHTableViewController {

    // MARK: - Private properties

    private var transaction: TransactionNotification

    // MARK: - Initialization

    init(transaction: TransactionNotification) {
        self.transaction = transaction

        super.init(title: "", items: [
            .notification(
                type: .verification,
                text: "This company is verified"),
            .description(
                date: transaction.date,
                title: "Personal details",
                description: "Please fill out your personal details."),
            .plugin,
            .form(placeholder: "First name"),
            .form(placeholder: "Last name"),
            .form(placeholder: "Date of birth"),
            .form(placeholder: "Address"),
            .form(placeholder: "Email"),
            .form(placeholder: "BSN number")
        ])
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Life cycle

    override func viewDidLoad() {
        super.viewDidLoad()

        tableView.allowsSelection = false
    }
}
