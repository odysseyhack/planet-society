//
//  TransactionLegalTermsViewController.swift
//  Permission Hub
//
//  Created by Corné on 14/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

final class TransactionLegalTermsViewController: PHTableViewController {

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
                image: nil,
                date: transaction.date,
                title: "Legal terms",
                description: "The Payment Services Directive 2 requires your express consent in order to allow the Third Party Provider (Stripe) to access your bank account information currently stored by your bank. This will allow Stripe to make payment on your behalf. You can withdraw your consent anytime."),
            .transactionItem(item: TransactionItem(
                item: "Privacy policy",
                fields: ["Read & accept the privacy policy"]),
                isChecked: false),
            .transactionItem(item: TransactionItem(
                item: "Terms & conditions",
                fields: ["Read & accept the terms & conditions"]),
                isChecked: false)
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
